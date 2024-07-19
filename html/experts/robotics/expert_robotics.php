<?php
require __DIR__."/settings.php";
require __DIR__."/vendor/Robot32lib/Middleware/Middleware.php";
require __DIR__."/vendor/Robot32lib/LLMServerList/LLMServerList.php";
require __DIR__."/vendor/Robot32lib/GPTlib/GPTlib.php";
require __DIR__."/vendor/Robot32lib/ULogger/ULogger.php";

require __DIR__."/vendor/Robot32lib/Biblio/Biblio.php";

use Robot32lib\LLMServerList\LLMServerList;
use Robot32lib\GPTlib\GPTlib;
use Robot32lib\ULogger\ULogger;
use Robot32lib\Biblio\Biblio;

ignore_user_abort(true); 
header('Content-Type:text/plain'); //NB. avoid xss
header('Meter-Bytes:true');

$LLMServerList = new LLMServerList();
$llms_to_try = $LLMServerList->getLLMFor("smart");

$ai = new GPTlib();

$bib = new Biblio(__DIR__."/bibliotheca");
$robotics_wisdom = trim(file_get_contents(__DIR__."/robotics_wisdom.txt"));
$robotics_wisdom.= "\n\n".$bib->getWisdom($_REQUEST["content"]);

$robotics_wisdom = trim(str_replace("\r","",$robotics_wisdom));
$robotics_wisdom = explode("\n\n",$robotics_wisdom); 
//print_r($robotics_wisdom);
$ai->setHistory($robotics_wisdom);

if($_REQUEST["history"]){
    $ai->setHistory($_REQUEST["history"],$ai->history);
}
   
$r = $ai->chat($_REQUEST["content"],$llms_to_try,null,function($txt,$data){
    if(!headers_sent() && $data)header("openrouter-id: ".$data['id']);
    echo $txt;
    @flush(); @ob_flush(); @ob_clean();
}); 

echo "\nRobotics answered!\n";

$logger = new Robot32lib\ULogger\ULogger($BASE_DIR);
$logger->log($_REQUEST["content"],"",$r['text'],$r['data']['id'],$r['data']['usage']['prompt_tokens'],$r['data']['usage']['completion_tokens'],$r['cost']);