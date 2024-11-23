{
  description = "A Nix flake for a Go project using buildGoModule";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs"; # Use Nixpkgs for dependencies
  };

  outputs = { self, nixpkgs }: {
    packages = nixpkgs.lib.genAttrs [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ] (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      pkgs.buildGoModule rec {
        pname = "projecthelper";
        version = "0.5.0";    # Replace with your version

        # The source code for the Go project
        src = ./.;

        installPhase = ''
          mkdir -p $out/bin
          cp $GOPATH/bin/projecthelper $out/bin/
        '';

        vendorHash = "sha256-Xq61Ji6gMP6OFdHSQqfrJVsKdatYM0ZezMbdz4Adr1A="; # nixpkgs.lib.fakeHash;

        # Metadata
        meta = with pkgs.lib; {
          description = "A Go project packaged with buildGoModule";
          license = licenses.mit; # Replace with your project's license
          maintainers = [ maintainers.yourGitHubHandle ];
          platforms = platforms.all; # Declare all platforms as supported
        };
      }
    );
  } ;
}

