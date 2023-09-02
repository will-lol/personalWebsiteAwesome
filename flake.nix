{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system: 
      let 
        name = "personalWebsiteCdk";
	overlays = [];
	lib = nixpkgs.lib;
	pkgs = import nixpkgs { inherit system overlays; };
	goBuild = (pkgs.buildGoModule {
	  inherit name;
	  src = ./.;
	  CGO_ENABLED = 0;
	  ldflags = ["-s -w"];
	  vendorSha256 = "sha256-RPa2SwK3YeUpOYQas0Y9+rprC3RNZKT2gGgplOMHNsk=";
	}).overrideAttrs (old: old // {GOOS = "linux"; GOARCH = "arm64"; });
      in
        {
	  packages = {
	    go = goBuild;
	  };
	  defaultPackage = goBuild;
	  devShell = pkgs.mkShell {
	    packages = [ pkgs.go-task pkgs.nodejs_18 pkgs.awscli2 pkgs.aws-sam-cli pkgs.go pkgs.gopls ];
	    shellHook = ''
              export GOPATH=$(pwd)/.go
	      export PATH=$GOPATH/bin:$PATH
	      export ENVIRONMENT=dev
	    '';
	  };
	}
    );
}
