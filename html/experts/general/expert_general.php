<?php
require __DIR__."/settings.php";
require __DIR__."/vendor/Robot32lib/Middleware/Middleware.php";
require __DIR__."/vendor/Robot32lib/LLMServerList/LLMServerList.php";
require __DIR__."/vendor/Robot32lib/GPTlib/GPTlib.php";
require __DIR__."/vendor/Robot32lib/ULogger/ULogger.php";
use Robot32lib\LLMServerList\LLMServerList;
use Robot32lib\GPTlib\GPTlib;
use Robot32lib\ULogger\ULogger;

ignore_user_abort(true); 
header('Content-Type:text/plain'); //NB. avoid xss
header('Meter-Bytes:true');

$LLMServerList = new LLMServerList();
$llms_to_try = $LLMServerList->getLoginFor("smart");

$ai = new GPTlib();
$ai->setHistory($_REQUEST["history"] ?? null);
$ai->setOptions([   
    "temperature"=> 1,  "max_tokens"=> 8024,
    "top_p"=> 1,        "stream"=> true,
    "stop"=> null                       
]);    

/*array_unshift($llms_to_try, ['url'=>'jjjj','headers'=>[
    "Authorization: Bearer yyyy",
    "Content-Type: application/json"
]]);*/

$r = $ai->chat($_REQUEST["content"],$llms_to_try,null,function($txt,$data){
    if(!headers_sent() && $data)header("openrouter-id: ".$data['id']);
    echo $txt;
    @flush(); @ob_flush(); @ob_clean();
}); 

$logger = new Robot32lib\ULogger\ULogger($BASE_DIR);
$logger->log($_REQUEST["content"],"",$r['text'],$r['data']['id'],$r['data']['usage']['prompt_tokens'],$r['data']['usage']['completion_tokens'],$r['cost']);


//llmcredentialsmanager
//$OPENROUTER_API_KEY = trim(file_get_contents(__DIR__."/../../../keys/openrouter.txt"));       
/*$headers = [
    "Authorization: Bearer $OPENROUTER_API_KEY",
    "Content-Type: application/json"
];
$url = "https://openrouter.ai/api/v1/chat/completions"; 
*/

//$logins_to_try = LLMServerList::getLoginFor("fast");



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


