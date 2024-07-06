<?php
ignore_user_abort(true); //NB, otherwise might skip billing
require __DIR__."/settings.php";
require __DIR__."/vendor/Robot32lib/Middleware/Middleware.php";

header('Content-Type:text/plain'); //NB. avoid xss
header('Meter-Bytes:true');

require __DIR__."/vendor/Robot32lib/GPTlib/GPTlib.php";

if(!isset($_COOKIE['r_user_key'])){
    echo "You need to connect your Openrouter account (go to SETTINGS). \n\nIt is an more expensive model (well relativly, its still cents or franction of cents but we cannot do it for free). This way you are paying these cents - not us :). ";
    exit();
}

//$OPENROUTER_API_KEY = trim(file_get_contents(__DIR__."/../../../keys/openrouter.txt"));       

$OPENROUTER_API_KEY = $_COOKIE['r_user_key'];

if(!ctype_alnum(str_replace("-","",$OPENROUTER_API_KEY))){
    echo "Provided key is not consisting of alphanumeric characters and '-'. ";
    exit();
}



$headers = [
    "Authorization: Bearer $OPENROUTER_API_KEY",
    "Content-Type: application/json"
];
$url = "https://openrouter.ai/api/v1/chat/completions"; 


$ai = new Robot32lib\GPTlib\GPTlib($url,$headers,TRUE);
$ai->setHistory($_REQUEST["history"] ?? null);

$options = [
    "temperature"=> 1,
    "max_tokens"=> 8024,
    "top_p"=> 1,
    "stream"=> true,
    "stop"=> null
    ];
$content = $_REQUEST["content"];
$model = $_REQUEST["model"];

$r = $ai->chat($content,$model,$options,function($txt,$data){
    if(!headers_sent())header("openrouter-id: ".$data['id']);
    
    //print_r($data);
    
    echo $txt;
    @flush(); @ob_flush(); @ob_clean();
}); 


if($r['error_code']){
    echo "Error ".$r['error_code']." ".$r['error'];
}
//print_r($r);
//$r = $ai->chat($_REQUEST["content"],$_REQUEST["model"],$history);




require __DIR__."/vendor/Robot32lib/ULogger/ULogger.php";

$logger = new Robot32lib\ULogger\ULogger($BASE_DIR);
$logger->log($content,$model,$r['text'],$r['data']['id'],$r['data']['usage']['prompt_tokens'],$r['data']['usage']['completion_tokens'],$r['cost']);


// ********************************** LOG IT *****************************
/*$currentTime = time();
$year = date('Y', $currentTime);
$month = date('m', $currentTime);
$day = date('d', $currentTime);
$hour = date('H', $currentTime);
$minute = date('i', $currentTime);
$second = date('s', $currentTime);
$time = "$year.$month.$day..$hour.$minute.$second";

$filename = $time . "___" . microtime(true);
$filename = str_replace(".", "_", $filename); // replace the decimal with an underscore
file_put_contents(__DIR__."/../../../collected_data/chats/".$filename.".txt", "Model: $model\n\n"."Query:\n".$content."\n\n\nResult:\n".$r['text']."\n\nCost:".$r['cost']);
*/


