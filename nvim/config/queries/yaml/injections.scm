;; Re-parse quoted and block scalars as yaml
(block_node
  (block_scalar
    (comment) @lang
  ) @injection.content
  (#offset! @injection.content 0 1 0 0) ; Prevent infinite recursion
  (#set-lang-from-hint! @lang) ; Use directive in ~/.vim/nvim.lua to trim comment into a filetype and set the injection language
)
