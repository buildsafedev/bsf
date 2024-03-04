
	{
	   lib,
	   stdenv,
	   buildGoModule,
	   ... 
	 }: buildGoModule {
	   name = "bsf";
	   src = ../.;  
	   doCheck = false;
	   
		vendorHash = "sha256-x/9O7qBzA1PaHBmZEzd8Kt0XuSw6feuJnd9FVUbdTt4=";
		
	   meta = with lib; {
		 description = "";
	   };
	   
	   
	 }
	