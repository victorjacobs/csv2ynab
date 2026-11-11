{
  lib,
  buildGo127Module,
}:

buildGo127Module {
  pname = "csv2ynab";
  version = "0.0.1";

  src = lib.cleanSource ../.;

  vendorHash = "sha256-BfeuBLKKZg9jVcAg2uLlMrmsVfT3xUelUTOTH5NLgyY=";

  meta = {
    description = "Import Belgian bank exports into YNAB";
    homepage = "https://github.com/victorjacobs/csv2ynab";
    license = lib.licenses.mit;
    mainProgram = "csv2ynab";
  };
}
