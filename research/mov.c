static int cbcs_scheme_decrypt(MOVContext *c, MOVStreamContext *sc, AVEncryptionInfo *sample, uint8_t *input, int size)
{
    int i, ret, rem_bytes;
    uint8_t iv[16];
    uint8_t *data;

    if (!sc->cenc.aes_ctx) {
        /* initialize the cipher */
        sc->cenc.aes_ctx = av_aes_alloc();
        if (!sc->cenc.aes_ctx) {
            return AVERROR(ENOMEM);
        }

        ret = av_aes_init(sc->cenc.aes_ctx, c->decryption_key, 16 * 8, 1);
        if (ret < 0) {
            return ret;
        }
    }

    /* whole-block full sample encryption */
    if (!sample->subsample_count) {
        /* decrypt the whole packet */
        memcpy(iv, sample->iv, 16);
        av_aes_crypt(sc->cenc.aes_ctx, input, input, size/16, iv, 1);
        return 0;
    } else if (!sample->crypt_byte_block && !sample->skip_byte_block) {
        av_log(c->fc, AV_LOG_ERROR, "pattern encryption is not present in 'cbcs' scheme\n");
        return AVERROR_INVALIDDATA;
    }

    for (i = 0; i < sample->subsample_count; i++) {
        if (sample->subsamples[i].bytes_of_clear_data + sample->subsamples[i].bytes_of_protected_data > size) {
            av_log(c->fc, AV_LOG_ERROR, "subsample size exceeds the packet size left\n");
            return AVERROR_INVALIDDATA;
        }

        /* skip the clear bytes */
        input += sample->subsamples[i].bytes_of_clear_data;
        size -= sample->subsamples[i].bytes_of_clear_data;

        /* decrypt the encrypted bytes */
        memcpy(iv, sample->iv, 16);
        data = input;
        rem_bytes = sample->subsamples[i].bytes_of_protected_data;
        while (rem_bytes > 0) {
            if (rem_bytes < 16*sample->crypt_byte_block) {
                break;
            }
            av_aes_crypt(sc->cenc.aes_ctx, data, data, sample->crypt_byte_block, iv, 1);
            data += 16*sample->crypt_byte_block;
            rem_bytes -= 16*sample->crypt_byte_block;
            data += FFMIN(16*sample->skip_byte_block, rem_bytes);
            rem_bytes -= FFMIN(16*sample->skip_byte_block, rem_bytes);
        }
        input += sample->subsamples[i].bytes_of_protected_data;
        size -= sample->subsamples[i].bytes_of_protected_data;
    }

    if (size > 0) {
        av_log(c->fc, AV_LOG_ERROR, "leftover packet bytes after subsample processing\n");
        return AVERROR_INVALIDDATA;
    }

    return 0;
}
