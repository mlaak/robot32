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

$r = $ai->chat($_REQUEST["content"],$llms_to_try,null,function($txt,$data){
    if(!headers_sent() && $data)header("openrouter-id: ".$data['id']);
    echo $txt;
    @flush(); @ob_flush(); @ob_clean();
}); 

$logger = new Robot32lib\ULogger\ULogger($BASE_DIR);
$logger->log($_REQUEST["content"],"",$r['text'],$r['data']['id'],$r['data']['usage']['prompt_tokens'],$r['data']['usage']['completion_tokens'],$r['cost']);

