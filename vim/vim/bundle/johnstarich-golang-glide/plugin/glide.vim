"command! GlideClearCache :exe '!glide cc'
command! GlideInstall :exe '!glide install -v'
command! GlideUpdate :exe '!glide update -v'
command! GlideUpdateAndInstall :exe '!glide update -v && glide install -v'
command! GlideQuickUpdate :exec '!glide update --no-recursive'
command! GlideQuickUpdateAndInstall :exec '!glide update --no-recursive -v && glide install -v'