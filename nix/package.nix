{
  lib,
  buildGo126Module,
}:

buildGo126Module {
  pname = "csv2ynab";
  version = "0.0.1";

  src = lib.cleanSource ../.;

  vendorHash = "sha256-exXUfKjWOoZy/4FFLKR+LJ544EdEAkx4cDoagcHg0aw=";

  meta = {
    description = "Import Belgian bank exports into YNAB";
    homepage = "https://github.com/victorjacobs/csv2ynab";
    license = lib.licenses.mit;
    mainProgram = "csv2ynab";
  };
}
