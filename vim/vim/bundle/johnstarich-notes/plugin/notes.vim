
command! -nargs=+ Notes :exec ":tabe " . fnameescape(system("notes " . <q-args>)) . "|cd %:p:h"
