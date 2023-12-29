
	{
	   lib,
	   stdenv,
	   buildGoModule,
	   ... 
	 }: buildGoModule {
	   name = "bsf";
	   src = ../.;  
	   
		vendorHash = "sha256-f29THF+FPQ4ORx2SJ2EJVOGbvoEvF4V+V0ZkmjUY35o=";
		
	   meta = with lib; {
		 description = "";
	   };
	 }
	