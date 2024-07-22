<?php
include __DIR__."/../../settings.php";

if(!function_exists("str_contains")){ //pre PHP 8
    function str_contains($haystack,$needle){
        return strpos($haystack,$needle)!==false;
    }
}
function run($cmd){
    $output = [];
    exec($cmd,$output);
    return strtolower(implode("\n",$output));
}
$G = "\033[1;32m"; //green text
$R = "\033[1;31m"; //red text
$N = "\033[0m"; //neutral - resets text

chdir(__DIR__."/../..");


//normally the user needs to bring their own, but we use ours for testing
$key = trim(file_get_contents($BASE_DIR."/keys/openrouter.txt")); 
$run = 'php expert_remote.php model="openai/gpt-4o-mini" r_user_key="'.$key.'" ';


//testing simple logic
$r = run($run.' content="What is 5 plus 11?"');

if(!str_contains($r,"16") && str_contains($r,"sixteen")){
    echo "$R REMOTE EXPERT TEST FAIL1 (simple logic):$N ";
    echo "Expected to get answer 16 but got".$r."\n";
} else {
    echo "$G pass remote expert test 1 (simple logic)$N\n";
}

//testing history 
$r = run($run.' history="user:pic a number;;ai:13" content="Multiply that with 2"');

if(!str_contains($r,"26") && !str_contains($r,"twenty six")){
    echo "$R REMOTE EXPERT TEST FAIL2 (chat history):$N ";
    echo "Expected to get answer 26 but got".$r."\n";
} else {
    echo "$G pass remote expert test 2 (chat history)$N\n";
}

//testing that we are talking to OpenAI 
$r = run($run.' content="Who are you?"');

if(!str_contains($r,"openai")){
    echo "$R REMOTE EXPERT TEST FAIL3 (openai):$N ";
    echo "Expected to get answer openai but got".$r."\n";
} else {
    echo "$G pass remote expert test 3 (openai)$N\n";
}
