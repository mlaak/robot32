<?php
chdir(__DIR__."/../..");

$G = "\033[1;32m"; //green text
$R = "\033[1;31m"; //red text
$N = "\033[0m"; //neutral - resets text

//testing simple logic
exec('php expert_general.php content="What is 5 plus 11?"',$output);
$result = strtolower(implode("\n",$output)); $output = [];

if(strpos($result,"16")===false && strpos($result,"sixteen")===false){
    echo "$R GENERAL EXPERT TEST FAIL1 (simple logic):$N ";
    echo "Expected to get answer 16 but got".$result."\n";
} else {
    echo "$G pass general expert test 1 (simple logic)$N\n";
}

//testing history 
exec('php expert_general.php history="user:pic a number;;ai:13" content="Multiply that with 2"',$output);
$result = strtolower(implode("\n",$output)); $output = [];

if(strpos($result,"26")===false && strpos($result,"twenty six")===false){
    echo "$R GENERAL EXPERT TEST FAIL2 (chat history):$N ";
    echo "Expected to get answer 26 but got".$result."\n";
} else {
    echo "$G pass general expert test 2 (chat history)$N\n";
}


