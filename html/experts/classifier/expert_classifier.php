<?php
require __DIR__."/settings.php";
require __DIR__."/vendor/Robot32lib/Middleware/Middleware.php";
require __DIR__."/vendor/Robot32lib/GPTlib/GPTlib.php";
require __DIR__."/vendor/Robot32lib/LLMServerList/LLMServerList.php";

ignore_user_abort(true); 
header('Content-Type:text/plain'); //NB. avoid xss
header('Meter-Bytes:true');

$LLMServerList = new Robot32lib\LLMServerList\LLMServerList();
$logins_to_try = $LLMServerList->getLoginFor("fast");
$ai = new Robot32lib\GPTlib\GPTlib($logins_to_try);

$ai->setOptions([
    "temperature"=> 1,  "max_tokens"=> 20,
    "top_p"=> 1,        "stream"=> false,
    "stop"=> null,      "curl_timeout" => 10,
    "curl_connect_timeout"=>5
]);

$user_query = str_replace('"','',$_REQUEST["content"]);
$classifier_text = file_get_contents(__DIR__."/classifier_text.txt");
$llm_query = $classifier_text."\n".'"'.$user_query.'"';

$result = $ai->chat($llm_query);
echo $result["text"];
