set completeopt=menu,menuone,noselect


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
"set statusline+=%{SyntasticStatuslineFlag()}
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

let g:syntastic_go_checkers = ['golangci-lint']
"let g:syntastic_go_gometalinter_args = ['--disable-all', '--enable=errcheck']
let g:syntastic_mode_map = { 'mode': 'active', 'passive_filetypes': ['go'] }

let g:syntastic_sh_shellcheck_args = "-x"

" UndoTree
map <leader>u <C-O>:UndotreeToggle<cr>


" Vim Go
let g:go_auto_type_info = 1
let g:go_highlight_extra_types = 1
let g:go_highlight_fields = 1
let g:go_highlight_format_strings = 1
let g:go_highlight_functions = 1
let g:go_highlight_function_calls = 1
let g:go_highlight_methods = 1
let g:go_highlight_operators = 1
let g:go_highlight_types = 1

" Goyo
autocmd! User GoyoEnter nested let g:goyo_previous_background = &background
autocmd! User GoyoLeave nested call SetBackground(g:goyo_previous_background)
" Pad resize fix from https://github.com/junegunn/goyo.vim/pull/104
"   technically, feedkeys is bad and should be "wincmd =", but that doesn't work
autocmd VimResized * call feedkeys("\<C-w>=")


" NERDTree
" Auto-open when vim opens
" autocmd vimenter * NERDTree
map <leader>t :NERDTreeToggle<CR>

" NERDTreeTabs
" Synchronize NERDTree with every tab
map <Leader>n <plug>NERDTreeTabsToggle<CR>

" EasyTags
"let g:easytags_cmd = '/usr/local/bin/ctags'
let g:easytags_async = 1

" vim-dotenv
" Automatically load env files in all parent directories at startup.
let s:dotenv_walked = 0
function! s:dotenv_walk_once()
  if s:dotenv_walked == 0
    call s:dotenv_walk()
    let s:dotenv_walked = 1
  endif
endfunction
function! s:dotenv_walk() abort
  let walk = findfile(".env", ".;", -1)
  for i in range(len(walk)-1, 0, -1)
      echo "dotenving " . walk[i]
      execute 'Dotenv ' . walk[i]
  endfor
endfunction
autocmd BufReadPre * call s:dotenv_walk_once()

" always show signcolumns (indicates problems on line)
set signcolumn=yes
" Make the sign column match the background
highlight SignColumn ctermbg=NONE
