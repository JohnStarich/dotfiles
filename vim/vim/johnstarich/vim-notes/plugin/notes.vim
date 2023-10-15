
command! -nargs=+ Notes :exec ":e " . fnameescape(system("notes " . <q-args> . " 2>/dev/null")) . "|cd %:p:h:h"
