<?php
require __DIR__."/settings.php";
require __DIR__."/vendor/Robot32lib/Middleware/Middleware.php";
require __DIR__."/vendor/Robot32lib/GPTlib/GPTlib.php";
require __DIR__."/vendor/Robot32lib/LLMServerList/LLMServerList.php";
require __DIR__."/vendor/Robot32lib/ClassTree/ClassTree.php";
use Robot32lib\ClassTree\ClassTree;
use Robot32lib\LLMServerList\LLMServerList;
use Robot32lib\GPTlib\GPTlib;

ignore_user_abort(true); 
header('Content-Type:text/plain'); //NB. avoid xss
header('Meter-Bytes:true');

$usrq = $_REQUEST["content"];

function removeCharacters($string, $charactersToRemove) {
    // Convert the characters to remove into an array if it's not already
    if (!is_array($charactersToRemove)) {
        $charactersToRemove = str_split($charactersToRemove);
    }   
    // Create an array of empty strings with the same length as $charactersToRemove
    $replacements = array_fill(0, count($charactersToRemove), '');  
    // Use str_replace to remove the specified characters
    return str_replace($charactersToRemove, $replacements, $string);
}

function pmatch($str,$pattern,$antipatterns = [],$remove = "!.?\n\r"){
    $str = removeCharacters($str,$remove);
    $str = strtolower($str);
    return preg_match($pattern, $str); 
}

$ct = new ClassTree();

$ptrn_joke = '/^(can\s)?(you\s)?(please\s)?tell\s(me\s)?(a\s)?joke(please\s)?$/i';
$ptrn_php = '/\bphp\b/i';
$ptrn_golang = '/\bgolang\b/i';
$ptrn_python = '/\bpython\b/i';
$ptrn_micropython = '/\bmicropython\b/i';
$ptrn_pi_pico = '/\\bpi pico\\b/i';

if(pmatch($usrq,$ptrn_joke))       {echo ClassTree::JOKE;       exit();}
if(pmatch($usrq,$ptrn_php))        {echo ClassTree::PHP;        exit();}
if(pmatch($usrq,$ptrn_golang))     {echo ClassTree::GOLANG;     exit();}
if(pmatch($usrq,$ptrn_python))     {echo ClassTree::PYTHON;     exit();}
if(pmatch($usrq,$ptrn_micropython)){echo ClassTree::MICROPYTHON;exit();}
if(pmatch($usrq,$ptrn_pi_pico))    {echo ClassTree::MICROCONTROLLERS;exit();}


$LLMServerList = new LLMServerList();
$llms_to_try = $LLMServerList->getLLMFor("fast");

$ai = new GPTlib();


$ai->setOptions([
    "temperature"=> 1,  "max_tokens"=> 20,
    "top_p"=> 1,        "stream"=> false,
    "stop"=> null,      "curl_timeout" => 10,
    "curl_connect_timeout"=>5
]);


$user_query = str_replace('"','',$_REQUEST["content"]);

$classifier_text = "Please help me to classify user request.The classifications are:\n"
                    .$ct->getTreeText()
                    ."Please answer with the classification code that best fits."
                    ."Please start your answer with the classification code.\n"
                    ."The user query is:";
//$classifier_text = file_get_contents(__DIR__."/classifier_text.txt");
$llm_query = $classifier_text."\n".'"'.$user_query.'"';


//echo $llm_query;

$result = $ai->chat($llm_query,$llms_to_try);
echo $result["text"];
