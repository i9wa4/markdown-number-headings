{
  description = "markdown-formatter - Markdown heading and table formatter";

  nixConfig = {
    extra-substituters = [
      "https://nix-community.cachix.org"
      "https://cache.numtide.com"
    ];
    extra-trusted-public-keys = [
      "nix-community.cachix.org-1:mB9FSh9qf2dCimDSUo8Zy7bkq5CX+/rkCWyvRCYg3Fs="
      "niks3.numtide.com-1:DTx8wZduET09hRmMtKdQDxNNthLQETkc/yaX7M4qK0g="
    ];
  };

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
    git-hooks = {
      url = "github:cachix/git-hooks.nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    inputs@{
      self,
      flake-parts,
      ...
    }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = [
        "aarch64-darwin"
        "x86_64-darwin"
        "x86_64-linux"
        "aarch64-linux"
      ];

      imports = [
        inputs.git-hooks.flakeModule
        inputs.treefmt-nix.flakeModule
      ];

      perSystem =
        {
          config,
          pkgs,
          self',
          ...
        }:
        let
          ghWorkflowFiles = "^\\.github/workflows/.*\\.(yml|yaml)$";
          go = pkgs.go_1_26;
          buildGoModule = pkgs.buildGoModule.override { inherit go; };
          version =
            if (builtins.hasAttr "ref" self && builtins.match "v[0-9]+\\.[0-9]+\\.[0-9]+" self.ref != null) then
              self.ref
            else if (builtins.hasAttr "shortRev" self) then
              "git-${self.shortRev}"
            else
              "dev";
          commit = if (builtins.hasAttr "rev" self) then builtins.substring 0 7 self.rev else "unknown";
          rumdlConfig = pkgs.writeText "rumdl.toml" ''
            disable = ["MD041"]

            [MD013]
            code-blocks = false
            headings = false
            reflow = true
          '';
          zizmorConfig = pkgs.writeText "zizmor.yml" ''
            rules:
              cache-poisoning:
                ignore:
                  - release.yml
          '';
        in
        {
          packages.default = buildGoModule {
            pname = "markdown-formatter";
            inherit version;
            src = ./.;
            vendorHash = null;
            ldflags = [
              "-s"
              "-w"
              "-X github.com/i9wa4/markdown-formatter/internal/version.Version=${version}"
              "-X github.com/i9wa4/markdown-formatter/internal/version.Commit=${commit}"
            ];
          };

          checks.build = self'.packages.default;

          devShells = {
            default = pkgs.mkShell {
              buildInputs = with pkgs; [
                actionlint
                deadnix
                gh
                ghalint
                gitleaks
                go
                gofumpt
                gopls
                govulncheck
                pinact
                rumdl
                statix
                treefmt
                zizmor
              ];
              shellHook = ''
                ${config.pre-commit.installationScript}
              '';
            };
            ci = pkgs.mkShell {
              buildInputs = with pkgs; [
                gitleaks
                go
                govulncheck
              ];
            };
            cd = pkgs.mkShell {
              buildInputs = with pkgs; [
                goreleaser
              ];
            };
          };

          apps = {
            check = {
              type = "app";
              program = "${pkgs.writeShellScriptBin "check" ''
                set -euo pipefail
                exec ${pkgs.nix}/bin/nix flake check --print-build-logs "$@"
              ''}/bin/check";
              meta.description = "Run flake checks.";
            };
          };

          treefmt = {
            projectRootFile = "flake.nix";
            programs = {
              nixfmt.enable = true;
              gofumpt.enable = true;
              shfmt = {
                enable = true;
                indent_size = 2;
              };
            };
            settings = {
              formatter = {
                rumdl = {
                  command = "${pkgs.rumdl}/bin/rumdl";
                  options = [
                    "fmt"
                    "--config"
                    "${rumdlConfig}"
                  ];
                  includes = [ "*.md" ];
                };
              };
              global.excludes = [
                ".direnv"
                ".git"
                "*.lock"
              ];
            };
          };

          pre-commit = {
            check.enable = true;
            settings.hooks = {
              end-of-file-fixer.enable = true;
              trim-trailing-whitespace.enable = true;
              check-added-large-files.enable = true;
              detect-private-keys.enable = true;
              check-merge-conflicts.enable = true;
              check-yaml.enable = true;
              actionlint.enable = true;
              ghalint = {
                enable = true;
                entry = "${pkgs.ghalint}/bin/ghalint run";
                files = ghWorkflowFiles;
              };
              pinact = {
                enable = true;
                entry = "${pkgs.pinact}/bin/pinact run";
                files = ghWorkflowFiles;
              };
              zizmor = {
                enable = true;
                entry = "${pkgs.zizmor}/bin/zizmor --config ${zizmorConfig}";
                files = ghWorkflowFiles;
              };
              statix = {
                enable = true;
                entry = "${pkgs.bash}/bin/bash -c '${pkgs.statix}/bin/statix check flake.nix'";
                pass_filenames = false;
              };
              deadnix.enable = true;
              rumdl-check = {
                enable = true;
                entry = "${pkgs.rumdl}/bin/rumdl check --config ${rumdlConfig}";
                types = [ "markdown" ];
              };
              govet = {
                enable = true;
                entry = "${pkgs.bash}/bin/bash -c 'test -n \"$NIX_BUILD_TOP\" || ${go}/bin/go vet ./...'";
                pass_filenames = false;
                types = [ "go" ];
              };
            };
          };
        };
    };
}
