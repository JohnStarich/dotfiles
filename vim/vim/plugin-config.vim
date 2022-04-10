
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
let g:go_fold_enable = ['import']
let g:go_build_tags = 'integration'

"let g:go_metalinter_autosave = 1
" g:go_auto_sameids Has trouble handling key input while running
let g:go_auto_sameids = 0
let g:go_jump_to_error = 0

let g:go_info_mode = 'gopls'
let g:go_def_mode = 'gopls'

let g:go_auto_type_info = 1
let g:go_highlight_extra_types = 1
let g:go_highlight_fields = 1
let g:go_highlight_format_strings = 1
let g:go_highlight_functions = 1
let g:go_highlight_function_calls = 1
let g:go_highlight_methods = 1
let g:go_highlight_operators = 1
let g:go_highlight_types = 1

let g:go_fmt_options = {
  \ 'gofmt': '-s',
  \ 'goimports': '-local github.ibm.com',
  \ }

noremap <leader>i :GoImports<cr>
noremap <leader>a :GoAlternate<cr>
noremap <leader>r :call go#lsp#Exit()<cr>
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

" Coc
" Use tab for trigger completion with characters ahead and navigate.
" Use command ':verbose imap <tab>' to make sure tab is not mapped by other plugin.
function! s:check_back_space() abort
  let col = col('.') - 1
  return !col || getline('.')[col - 1]  =~# '\s'
endfunction
inoremap <silent><expr> <TAB>
      \ pumvisible() ? "\<C-n>" :
      \ <SID>check_back_space() ? "\<TAB>" :
      \ coc#refresh()
inoremap <expr><S-TAB> pumvisible() ? "\<C-p>" : "\<C-h>"
" reformat the current buffer
command! -nargs=0 Format :call CocAction('format')
" always show signcolumns (indicates problems on line)
set signcolumn=yes
" Make the sign column match the background
highlight SignColumn ctermbg=NONE
" Better display for messages at the bottom
set cmdheight=2

" Coc Spell Checker
" coc-spell-checker can use these, but multipurpose
vmap <leader>a <Plug>(coc-codeaction-selected)
nmap <leader>a <Plug>(coc-codeaction-selected)

" vim-dotenv
" Automatically load env files in all parent directories at startup.
function! s:dotenv_walk() abort
  let walk = findfile(".env", ".;", -1)
  for i in range(len(walk)-1, 0, -1)
      execute 'Dotenv ' . walk[i]
  endfor
endfunction
autocmd VimEnter * call s:dotenv_walk()
