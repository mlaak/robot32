<?php
ignore_user_abort(true); //NB, otherwise might skip billing
require __DIR__."/settings.php";
require __DIR__."/vendor/Robot32lib/Middleware/Middleware.php";

header('Content-Type:text/plain'); //NB. avoid xss
//header('Meter-Bytes:true');
srand();
$r = rand(0,3);

switch($r){
    case 0: echo "Who said PHP developers are poor? We see dollar signs all the time!";break;
    case 1: echo "Why developers confuse Christmas and Halloween? Because OCT 31 == DEC 25.";break;
    case 2: echo "How to understand recursion? For that, you need to first understand recursion.";break;
    case 3: echo "Why do submarines run Linux? Because you cant open windows under water.";break;
}


