{
  description = "Murtaza Udaipurwala's Go Monolith";
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-23.11";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = { self, nixpkgs, flake-utils, ... }@inputs:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        formatter = pkgs.nixpkgs-fmt;
        packages = {
          default = pkgs.buildGoModule {
            pname = "z";
            version = "0.2.1";
            src = ./.;
            vendorHash = "sha256-hDTZG7M5HdVZks7yya7dVDC6NFyEkY3KjGIdOcZX9Ro=";
            CGO_ENABLED = 0;
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
            go
            go-tools
            gopls
            nixpkgs-fmt
            nixd
          ];
        };
      });
}
