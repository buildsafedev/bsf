
	{
	   lib,
	   stdenv,
	   buildGoModule,
	   ... 
	 }: buildGoModule {
	   name = "";
	   src = ../.;  
	   
		vendorHash = "sha256-f29THF&#43;FPQ4ORx2SJ2EJVOGbvoEvF4V&#43;V0ZkmjUY35o=";
		
	   meta = with lib; {
		 description = "";
	   };
	 }
	