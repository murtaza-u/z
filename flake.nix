{
  description = "Monolith Go Commander";
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        formatter = pkgs.nixpkgs-fmt;
        packages = {
          default = pkgs.buildGoModule {
            pname = "z";
            version = "0.1.1";
            src = ./.;
            vendorHash = "sha256-WlwDfQxojWlSu455ZH0KEbb6eOi/rm7eap6XCs6lqjU=";
            env.CGO_ENABLED = 0;
            subPackages = [ "cmd/z" ];
            nativeBuildInputs = [ pkgs.installShellFiles ];
            postInstall = ''
              for shell in bash zsh; do
                installShellCompletion --$shell ./completion/$shell/z
              done
            '';
          };
        };
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            nixpkgs-fmt
            nixd
            go
            go-tools
            gopls
          ];
        };
      });
}
