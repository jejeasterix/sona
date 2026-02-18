#pragma once
#include <whisper.h>
#include <stdint.h>

void sona_whisper_set_verbose(int verbose);
void sona_whisper_set_stream_callbacks(struct whisper_full_params *params, uintptr_t handle);
