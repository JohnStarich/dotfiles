
imap <F7> <Plug>(JavaComplete-Imports-RemoveUnused)


" SuperTab

" Use omnifunc for completion, always
"let g:SuperTabDefaultCompletionType = "<c-x><c-o>"


" Vim Markdown

" Disable Vim Markdown folding
let g:vim_markdown_folding_disabled = 1
" Enable syntax name aliases for fenced code blocks
"let g:vim_markdown_fenced_languages = ['csharp=cs']
" Highlight YAML front matter
let g:vim_markdown_frontmatter = 1
" Follow markdown links without .md extension with `ge`
let g:vim_markdown_no_extensions_in_markdown = 1
" Automatic bullet insertion
"let g:vim_markdown_new_list_item_indent = 0
"autocmd BufNewFile,BufRead markdown setlocal formatoptions-=r
"autocmd BufNewFile,BufRead markdown setlocal formatoptions-=c
"autocmd BufNewFile,BufRead markdown setlocal formatoptions-=o


" DelimitMate

let delimitMate_expand_cr = 2
let delimitMate_quotes = "\" ' `"
let delimitMate_nesting_quotes = ['"', '`']
" autocmd FileType markdown let b:delimitMate_quotes = \" ' ` _ *
autocmd FileType markdown let b:delimitMate_quotes = "\" ' ` _"
autocmd FileType markdown let b:delimitMate_nesting_quotes = ['`', '_', '*']
autocmd FileType python,markdown let b:delimitMate_expand_inside_quotes = 1


" Syntastic

let g:syntastic_shell = "/bin/bash"

set statusline+=%#warningmsg#
set statusline+=%{SyntasticStatuslineFlag()}
set statusline+=%*

let g:syntastic_always_populate_loc_list = 1
" show a list of failed checks in bottom pane
"let g:syntastic_auto_loc_list = 1
" run checker on open
"let g:syntastic_check_on_open = 1
" run checker on quit, 0 to disable
let g:syntastic_check_on_wq = 0
" show all errors
let g:syntastic_aggregate_errors = 1

let g:syntastic_python_checkers = ['flake8']

"let g:syntastic_java_maven_executable = '/usr/local/Cellar/maven/3.3.9/bin/mvn'

" UndoTree
map <leader>u <C-O>:UndotreeToggle<cr>
