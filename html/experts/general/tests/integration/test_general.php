<?php
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


//testing simple logic
$r = run('php expert_general.php content="What is 5 plus 11?"');

if(!str_contains($r,"16") && !str_contains($r,"sixteen")){
    echo "$R GENERAL EXPERT TEST FAIL1 (simple logic):$N ";
    echo "Expected to get answer 16 but got".$r."\n";
} else {
    echo "$G pass general expert test 1 (simple logic)$N\n";
}

//testing history 
$r=run('php expert_general.php history="user:pic a number;;ai:13" content="Multiply that with 2"');

if(!str_contains($r,"26") && !str_contains($r,"twenty six")){
    echo "$R GENERAL EXPERT TEST FAIL2 (chat history):$N ";
    echo "Expected to get answer 26 but got".$r."\n";
} else {
    echo "$G pass general expert test 2 (chat history)$N\n";
}


