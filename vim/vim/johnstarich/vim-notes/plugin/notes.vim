
command! -nargs=+ Notes :exec ":e " . fnameescape(system("notes " . <q-args>)) . "|cd %:p:h:h"
