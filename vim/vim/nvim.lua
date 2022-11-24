-- ray-x/go.nvim
require('go').setup {
    -- https://github.com/ray-x/go.nvim#configuration
    textobjects = false, -- Disables 'treesitter' error
    -- gocoverage_sign = "█",
    gocoverage_sign = " ",
}
vim.cmd([[
augroup my-go-nvim-coverage
    au Syntax go hi goCoverageCovered guibg=DarkGreen ctermbg=DarkGreen
    au Syntax go hi goCoverageUncover guibg=DarkRed   ctermbg=DarkRed
augroup end
]])

local goCoverage = require('go.coverage')
vim.keymap.set('n', '<leader>cr', goCoverage.run, {})
vim.keymap.set('n', '<leader>ct', goCoverage.toggle, {})

-- telescope.nvim
-- See here for ideas: https://github.com/nvim-telescope/telescope.nvim/wiki/Configuration-Recipes
local telescope = require('telescope')
local telescopeActions = require("telescope.actions")
telescope.setup {
    defaults = {
        mappings = {
            i = {
                ["<esc>"] = telescopeActions.close
            },
        },
        prompt_prefix = "🚀 ",
        selection_caret = "👉 ",
        entry_prefix = "   ",
    },
    pickers = {
        -- push_tagstack_on_edit = true, -- Not documented, could go away: https://github.com/nvim-telescope/telescope.nvim/pull/1887
        find_files = {
            theme = "dropdown",
            push_tagstack_on_edit = true,
        },
        live_grep = {
            push_tagstack_on_edit = true,
        },
        lsp_dynamic_workspace_symbols = {
            theme = "dropdown",
            push_tagstack_on_edit = true,
        },
    },
    extensions = {
        'fzf',
    },
}

local telescopeBuiltin = require('telescope.builtin')
vim.keymap.set('n', '<leader>fb', telescopeBuiltin.buffers, {})
vim.keymap.set('n', '<leader>ff', telescopeBuiltin.find_files, {})
vim.keymap.set('n', '<leader>fg', telescopeBuiltin.live_grep, {})
vim.keymap.set('n', '<leader>fh', telescopeBuiltin.help_tags, {})
vim.keymap.set('n', '<leader>fs', telescopeBuiltin.lsp_dynamic_workspace_symbols, {})
vim.keymap.set('n', '<leader>fn', function()
    local notesDir = vim.env.NOTES_BASE or "~/notes"
    return telescopeBuiltin.live_grep {
        prompt_title = "Search notes",
        cwd = notesDir,
        search_dirs = { notesDir },
        additional_args = function(opts)
            return { "--follow" }
        end,
    }
end, {})
vim.keymap.set('n', '<leader>fv', function()
    return telescopeBuiltin.live_grep {
        prompt_title = "Search vim plugins",
        cwd = "~/.vim/bundle",
        search_dirs = { "~/.vim/bundle" },
    }
end, {})

-- nvim-lspconfig:
-- Mappings.
-- See `:help vim.diagnostic.*` for documentation on any of the below functions
local opts = { noremap=true, silent=true }
-- vim.keymap.set('n', '<space>e', vim.diagnostic.open_float, opts)
vim.keymap.set('n', '[d', vim.diagnostic.goto_prev, opts)
vim.keymap.set('n', ']d', vim.diagnostic.goto_next, opts)
-- vim.keymap.set('n', '<space>q', vim.diagnostic.setloclist, opts)

vim.keymap.set('n', '<leader>k', function()
    vim.lsp.stop_client(vim.lsp.get_active_clients())
end, opts)

-- Use an on_attach function to only map the following keys
-- after the language server attaches to the current buffer
local on_attach = function(client, bufnr)
    -- Enable completion triggered by <c-x><c-o>
    vim.api.nvim_buf_set_option(bufnr, 'omnifunc', 'v:lua.vim.lsp.omnifunc')

    -- Temporary work-around for telescope pickers using wrong cwd.
    vim.fn.chdir(".")

    -- Mappings.
    -- See `:help vim.lsp.*` for documentation on any of the below functions
    local bufopts = { noremap=true, silent=true, buffer=bufnr }
    -- Workspace commands.
    vim.keymap.set('n', '<space>wa', vim.lsp.buf.add_workspace_folder, bufopts)
    vim.keymap.set('n', '<space>wr', vim.lsp.buf.remove_workspace_folder, bufopts)
    vim.keymap.set('n', '<space>wl', function()
        print(vim.inspect(vim.lsp.buf.list_workspace_folders()))
    end, bufopts)
    -- Code jumps and docs.
    vim.keymap.set('n', 'K', vim.lsp.buf.hover, bufopts)
    vim.keymap.set('n', '<C-]>', vim.lsp.buf.definition, bufopts)
    vim.keymap.set('n', 'gi', vim.lsp.buf.implementation, bufopts)
    vim.keymap.set('n', 'gr', vim.lsp.buf.references, bufopts)
    -- Code actions.
    vim.keymap.set('n', '<leader>fo', vim.lsp.buf.formatting_sync, bufopts)
    vim.keymap.set('n', '<leader>r', vim.lsp.buf.rename, bufopts)
    vim.keymap.set('n', '<leader>ca', vim.lsp.buf.code_action, bufopts)

    -- Go mappings.
    -- Switch to alternate file. Code <-> Test.
    vim.keymap.set('n', '<leader>a', function()
        -- Switch to test file or vice versa:
        if vim.fn.expand('%:e') ~= "go" then
            return
        end
        local alternateGoFile = vim.fn.expand('%:r')
        local testSuffix = '_test'
        local testSuffixLen = string.len(testSuffix)
        if string.sub(alternateGoFile, -testSuffixLen) == testSuffix then
            alternateGoFile = string.sub(alternateGoFile, 1, string.len(alternateGoFile) - testSuffixLen)
        else
            alternateGoFile = alternateGoFile .. testSuffix
        end
        alternateGoFile = alternateGoFile .. '.' .. vim.fn.expand('%:e')
        vim.api.nvim_command("edit " .. alternateGoFile)
    end, bufopts)
    -- Run equivalent of gofmt and goimports on save.
    vim.api.nvim_create_autocmd("BufWritePre", {
        -- Format & goimports:
        pattern = { "*.go" },
        callback = function()
            vim.lsp.buf.formatting_sync()

            local params = vim.lsp.util.make_range_params(nil, "utf-16")
            params.context = { only = { "source.organizeImports" } }
            local result = vim.lsp.buf_request_sync(0, "textDocument/codeAction", params, 3000)
            for _, res in pairs(result or {}) do
                for _, r in pairs(res.result or {}) do
                    if r.edit then
                        vim.lsp.util.apply_workspace_edit(r.edit, "utf-16")
                    else
                        vim.lsp.buf.execute_command(r.command)
                    end
                end
            end
        end,
    })
end


-- Set up nvim-cmp for tab completions.
local cmp = require('cmp')
-- Current issues:
-- * Wrong completion for '%': https://github.com/hrsh7th/nvim-cmp/issues/1058

-- https://github.com/hrsh7th/nvim-cmp/wiki/Example-mappings#vim-vsnip
local has_words_before = function()
    local line, col = unpack(vim.api.nvim_win_get_cursor(0))
    return col ~= 0 and vim.api.nvim_buf_get_lines(0, line - 1, line, true)[1]:sub(col, col):match("%s") == nil
end

local feedkey = function(key, mode)
    vim.api.nvim_feedkeys(vim.api.nvim_replace_termcodes(key, true, true, true), mode, true)
end

cmp.setup({
    preselect = cmp.PreselectMode.None,
    snippet = {
        -- REQUIRED - you must specify a snippet engine
        expand = function(args)
            vim.fn["vsnip#anonymous"](args.body) -- For `vsnip` users.
        end,
    },
    window = {
        -- completion = cmp.config.window.bordered(),
        -- documentation = cmp.config.window.bordered(),
    },
    mapping = cmp.mapping.preset.insert({
    ['<C-d>'] = cmp.mapping.scroll_docs(-4),
    ['<C-f>'] = cmp.mapping.scroll_docs(4),
    ['<C-Space>'] = cmp.mapping.complete(),
    ['<CR>'] = cmp.mapping.confirm {
        behavior = cmp.ConfirmBehavior.Replace,
        select = true,
    },
    ['<Tab>'] = cmp.mapping(function(fallback)
            if cmp.visible() then
                cmp.select_next_item()
            elseif vim.fn["vsnip#available"](1) == 1 then
                feedkey("<Plug>(vsnip-expand-or-jump)", "")
            elseif has_words_before() then
                cmp.complete()
            else
                fallback()
            end
        end, { 'i', 's' }),
    ['<S-Tab>'] = cmp.mapping(function(fallback)
            if cmp.visible() then
                cmp.select_prev_item()
            elseif vim.fn["vsnip#jumpable"](-1) == 1 then
                feedkey("<Plug>(vsnip-jump-prev)", "")
            end
        end, { 'i', 's' }),
    }),
    sources = cmp.config.sources(
        -- Order determines completion ordering.
        {
            { name = 'nvim_lsp' },
            { name = 'nvim_lsp_signature_help' },
        },
        { name = 'vsnip' },
        { name = 'nvim_lua' },
        { name = 'buffer' },
        {})
})

-- Set configuration for specific filetype.
cmp.setup.filetype('gitcommit', {
    sources = cmp.config.sources({
        { name = 'cmp_git' }, -- You can specify the `cmp_git` source if you were installed it.
    }, {
        { name = 'buffer' },
    })
})

-- Use buffer source for `/` and `?` (if you enabled `native_menu`, this won't work anymore).
cmp.setup.cmdline({ '/', '?' }, {
    mapping = cmp.mapping.preset.cmdline(),
    sources = {
        { name = 'buffer' }
    }
})

-- Use cmdline & path source for ':' (if you enabled `native_menu`, this won't work anymore).
cmp.setup.cmdline(':', {
    mapping = cmp.mapping.preset.cmdline(),
    sources = cmp.config.sources({
        { name = 'path' }
    }, {
        { name = 'cmdline' }
    })
})

-- Set up lspconfig.
local servers = {
    gopls = {
        cmd = {'gopls', '-remote=auto' },
        settings = {
            gopls = {
                -- https://cs.opensource.google/go/x/tools/+/refs/tags/gopls/v0.10.1:gopls/doc/settings.md
                experimentalPostfixCompletions = true,
                analyses = {
                    unusedparams = true,
                    shadow = true,
                },
                staticcheck = true,
            },
        },
    },
}
local capabilities = require('cmp_nvim_lsp').default_capabilities()
local lspconfig = require('lspconfig')
for lsp, config in pairs(servers) do
    local opts = {
        on_attach = on_attach,
        capabilities = capabilities,
    }
    for key, value in pairs(config) do
        opts[key] = value
    end
    lspconfig[lsp].setup(opts)
end