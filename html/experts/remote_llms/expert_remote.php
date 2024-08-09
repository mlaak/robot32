<?php
require __DIR__."/settings.php";
require __DIR__."/vendor/Robot32lib/Middleware/Middleware.php";
require __DIR__."/vendor/autoload.php";
use Robot32lib\GPTlib\GPTlib;
use Robot32lib\ULogger\ULogger;
declare(ticks=1);

$arh = apache_request_headers();
if(isset($arh["Ttd"])){
    TTD_INIT($arh["Ttd"]*1);
}

/*function TDD_SET($c){
    global $TDD_c;
    $TDD_c = $c;
}

function TTD($message,...$params){*/

ignore_user_abort(true); 
header('Content-Type:text/plain'); //NB. avoid xss
header('Meter-Bytes:true');

if(isset($_REQUEST['r_user_key'])){
    $_COOKIE['r_user_key'] = $_REQUEST['r_user_key'];
}

if(!isset($_COOKIE['r_user_key'])){
    echo "You need to connect your Openrouter account (go to SETTINGS). \n\nIt is an more expensive model (well relativly, its still cents or franction of cents but we cannot do it for free). This way you are paying these cents - not us :). ";
    exit();
}



TTD("HERE!!!!!!!!!!!!!!!","varname","val");

//$OPENROUTER_API_KEY = trim(file_get_contents(__DIR__."/../../../keys/openrouter.txt"));       
$OPENROUTER_API_KEY = $_COOKIE['r_user_key']; //we require user have their own key (expensive models)

if(!ctype_alnum(str_replace("-","",$OPENROUTER_API_KEY))){
    echo "Provided key is not consisting of alphanumeric characters and '-'. ";
    exit();
}

$headers = [
    "Authorization: Bearer $OPENROUTER_API_KEY",
    "Content-Type: application/json"
];

$ai = new GPTlib($URL,$headers,TRUE);
$ai->setHistory(TTX($_REQUEST["history"] ?? null));

$options = [
    "temperature"=> 1,
    "max_tokens"=> 8000,
    "top_p"=> 1,
    "stream"=> true,
    "stop"=> null
    ];
$content = TTX($_REQUEST["content"]);
$model = TTX($_REQUEST["model"]);

$r = $ai->chat($content,$model,$options,function($txt,$data){
    if(!headers_sent() && isset($data['id'])){
        header("openrouter-id: ".TTX($data['id']));
    }
    echo TTX($txt); //send the piece of prompt response to the user
    @flush(); @ob_flush(); @ob_clean(); //make sure data gets to user ASAP
}); 

if($r['error_code']){
    echo "Error ".TTX($r['error_code'])." ".TTX($r['error']);
}

$logger = new ULogger($BASE_DIR);
$logger->log($content,$model,$r['text'],$r['data']['id'],$r['data']['usage']['prompt_tokens'],$r['data']['usage']['completion_tokens'],$r['cost']);