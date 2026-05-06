const fs = require('fs');
const https = require('https');
const nodemailer = require('nodemailer');

const STATE_FILE = '.sona-pr-state.json';
const REPO_OWNER = 'thewh1teagle';
const REPO_NAME = 'sona';
const PR_NUMBER = 17;
const GH_TOKEN = process.env.GH_TOKEN;
const GMAIL_USER = process.env.GMAIL_USER;
const GMAIL_PASSWORD = process.env.GMAIL_PASSWORD;
const EMAIL_TO = process.env.EMAIL_TO;

// Read previous state
function readState() {
  if (fs.existsSync(STATE_FILE)) {
    try {
      return JSON.parse(fs.readFileSync(STATE_FILE, 'utf8'));
    } catch (e) {
      return {};
    }
  }
  return {};
}

// Write state
function writeState(state) {
  fs.writeFileSync(STATE_FILE, JSON.stringify(state, null, 2));
}

// GitHub API call
function githubApi(path) {
  return new Promise((resolve, reject) => {
    const options = {
      hostname: 'api.github.com',
      path: `/repos/${REPO_OWNER}/${REPO_NAME}${path}`,
      method: 'GET',
      headers: {
        'Authorization': `token ${GH_TOKEN}`,
        'User-Agent': 'sona-pr-monitor'
      }
    };

    https.request(options, (res) => {
      let data = '';
      res.on('data', chunk => data += chunk);
      res.on('end', () => {
        try {
          resolve(JSON.parse(data));
        } catch (e) {
          reject(e);
        }
      });
    }).on('error', reject).end();
  });
}

// Send email
async function sendEmail(subject, body) {
  const transporter = nodemailer.createTransport({
    host: 'smtp.gmail.com',
    port: 465,
    secure: true,
    auth: {
      user: GMAIL_USER,
      pass: GMAIL_PASSWORD
    }
  });

  await transporter.sendMail({
    from: GMAIL_USER,
    to: EMAIL_TO,
    subject: subject,
    text: body
  });

  console.log('✅ Email sent');
}

// Check for release newer than v0.2.1 (2026-03-13)
async function checkForNewRelease() {
  const releases = await githubApi('/releases?per_page=3');
  const targetDate = new Date('2026-03-13').getTime();

  for (const release of releases) {
    const publishedAt = new Date(release.published_at).getTime();
    if (publishedAt > targetDate && !release.draft && !release.prerelease) {
      return {
        found: true,
        tag: release.tag_name,
        name: release.name,
        date: release.published_at
      };
    }
  }
  return { found: false };
}

// Check if release contains no_speech_prob
async function releaseHasNoSpeechProb(tag) {
  try {
    const content = await githubApi(`/contents/internal/server/format.go?ref=${tag}`);
    const decoded = Buffer.from(content.content, 'base64').toString('utf8');
    return decoded.includes('no_speech_prob');
  } catch (e) {
    console.log(`Could not check ${tag}: ${e.message}`);
    return false;
  }
}

async function main() {
  try {
    const previousState = readState();

    // Get PR info
    const pr = await githubApi(`/pulls/${PR_NUMBER}`);
    const prState = pr.state;
    const updatedAt = pr.updated_at;
    const reviewCount = pr.reviews ? pr.reviews.length : 0;
    const commentCount = pr.comments || 0;

    // Get releases
    const newRelease = await checkForNewRelease();
    let releaseContainsPatch = false;
    if (newRelease.found) {
      releaseContainsPatch = await releaseHasNoSpeechProb(newRelease.tag);
    }

    // Check for changes
    const stateChanged = previousState.prState !== prState;
    const hasNewReviews = previousState.reviewCount < reviewCount;
    const hasNewComments = previousState.commentCount < commentCount;
    const hasNewRelease = newRelease.found && !previousState.latestReleaseTag;
    const releaseNowContainsPatch = releaseContainsPatch && !previousState.releaseHasNoSpeechProb;

    const shouldEmail =
      (prState === 'CLOSED' && previousState.prState !== 'CLOSED') ||
      (prState === 'MERGED' && previousState.prState !== 'MERGED') ||
      hasNewReviews ||
      hasNewComments ||
      hasNewRelease ||
      releaseNowContainsPatch;

    if (shouldEmail) {
      let emailBody = '[VS Transcript] sona#17 a bouge -- action possible\n\n';

      if (prState === 'MERGED' && previousState.prState !== 'MERGED') {
        emailBody += 'La PR sona#17 a ete FUSIONNEE ! 🎉\n';
        emailBody += 'Action : Telecharger la nouvelle release et mettre a jour binaries/sona-x86_64-pc-windows-msvc.exe\n\n';
      } else if (prState === 'CLOSED' && previousState.prState !== 'CLOSED') {
        emailBody += 'La PR sona#17 a ete FERMEE.\n';
        emailBody += 'Action : Verifier pourquoi et relancer si necessaire.\n\n';
      }

      if (hasNewReviews) {
        emailBody += `- ${reviewCount} review(s) recu(e)s\n`;
        emailBody += 'Action : Relever les commentaires et repondre si necessaire.\n\n';
      }

      if (hasNewComments) {
        emailBody += `- ${commentCount} commentaire(s) total\n`;
      }

      if (newRelease.found) {
        emailBody += `- Nouvelle release: ${newRelease.tag} (${newRelease.date.split('T')[0]})\n`;
        if (releaseContainsPatch) {
          emailBody += '  ✅ Contient le champ no_speech_prob\n';
          emailBody += '  Action : Telecharger et mettre a jour le binaire sona.\n\n';
        }
      }

      emailBody += `Lien: https://github.com/thewh1teagle/sona/pull/${PR_NUMBER}\n`;

      await sendEmail('[VS Transcript] sona#17 a bouge -- action possible', emailBody);
    } else {
      console.log(`✓ PR sona#17 inchangee, derniere maj : ${updatedAt}`);
      console.log(`✓ Aucune nouvelle release sona depuis v0.2.1 (2026-03-13)`);
    }

    // Save state
    writeState({
      prState,
      updatedAt,
      reviewCount,
      commentCount,
      latestReleaseTag: newRelease.tag || previousState.latestReleaseTag,
      releaseHasNoSpeechProb: releaseContainsPatch || previousState.releaseHasNoSpeechProb,
      lastCheck: new Date().toISOString()
    });

  } catch (error) {
    console.error('❌ Error:', error.message);
    process.exit(1);
  }
}

main();
