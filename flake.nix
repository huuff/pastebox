{
  description = "A server-side-rendered app for sharing your pastes";

  outputs = { self, nixpkgs }: 
  let
    system = "x86_64-linux";
    pkgs = import nixpkgs { inherit system; };
  in
  {

    devShell.${system} = pkgs.mkShell {
      buildInputs = with pkgs; [
        go 
      ];
    };

  };
}
