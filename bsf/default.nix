
	{
	   lib,
	   stdenv,
	   buildGoModule,
	   ... 
	 }: buildGoModule {
	   name = "";
	   src = ../.;  
	   
		vendorHash = lib.fakeHash;
		
	   meta = with lib; {
		 description = "";
	   };
	 }
	