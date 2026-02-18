#include "whisper_cgo.h"
#include <stdio.h>

// Forward declarations for Go-exported callback trampolines.
extern void sonaGoProgressCB(uintptr_t handle, int32_t progress);
extern void sonaGoSegmentCB(uintptr_t handle, void *ctx_ptr, int32_t n_new);
extern int32_t sonaGoAbortCB(uintptr_t handle);

static int sona_whisper_verbose = 0;

static void sona_whisper_log_callback(enum ggml_log_level level, const char * text, void * user_data) {
    (void) level;
    (void) user_data;
    if (!sona_whisper_verbose) {
        return;
    }
    fputs(text, stderr);
}

void sona_whisper_set_verbose(int verbose) {
    sona_whisper_verbose = verbose;
    whisper_log_set(sona_whisper_log_callback, NULL);
}

static void sona_whisper_progress_trampoline(struct whisper_context *ctx, struct whisper_state *state, int progress, void *user_data) {
    (void)ctx; (void)state;
    sonaGoProgressCB((uintptr_t)user_data, (int32_t)progress);
}

static void sona_whisper_new_segment_trampoline(struct whisper_context *ctx, struct whisper_state *state, int n_new, void *user_data) {
    (void)state;
    sonaGoSegmentCB((uintptr_t)user_data, ctx, (int32_t)n_new);
}

static _Bool sona_whisper_abort_trampoline(void *user_data) {
    return sonaGoAbortCB((uintptr_t)user_data) != 0;
}

void sona_whisper_set_stream_callbacks(struct whisper_full_params *params, uintptr_t handle) {
    void *h = (void *)handle;
    params->progress_callback = sona_whisper_progress_trampoline;
    params->progress_callback_user_data = h;
    params->new_segment_callback = sona_whisper_new_segment_trampoline;
    params->new_segment_callback_user_data = h;
    params->abort_callback = sona_whisper_abort_trampoline;
    params->abort_callback_user_data = h;
}
