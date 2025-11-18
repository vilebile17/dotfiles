return {
  {
    "catppuccin/nvim",
    name = "catppuccin",
    priority = 1000,
    config = function()
      require("catppuccin").setup({
        flavour = "mocha",
        integrations = {
          lazy = true,
          treesitter = true,
          native_lsp = { enabled = true },
          telescope = true,
          which_key = true,
          noice = true,
          cmp = true,
          gitsigns = true,
        },
      })
      vim.cmd.colorscheme("catppuccin-mocha")
    end,
  },
}
