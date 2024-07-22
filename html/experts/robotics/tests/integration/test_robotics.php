<?php
chdir(__DIR__."/../..");

$G = "\033[1;32m"; //green text
$R = "\033[1;31m"; //red text
$N = "\033[0m"; //neutral - resets text

//testing simple logic
exec('php expert_robotics.php content="What is 5 plus 11?"',$output);
$result = implode("\n",$output); $output = [];
if(strpos($result,"16")===false){
    echo "$R ROBOTICS EXPERT TEST FAIL1 (simple logic):$N ";
    echo "Expected to get answer 16 but got".$result."\n";
} else {
    echo "$G pass robotics expert test 1 (simple logic)$N\n";
}

//testing history 
exec('php expert_robotics.php history="user:pic a number;;ai:13" content="Multiply that with 2"',$output);
$result = implode("\n",$output); $output = [];
if(strpos($result,"26")===false){
    echo "$R ROBOTICS EXPERT TEST FAIL2 (chat history):$N ";
    echo "Expected to get answer 26 but got".$result."\n";
} else {
    echo "$G pass robotics expert test 2 (chat history)$N\n";
}


//testing that it would recommend pi pico microcontroller
exec('php expert_robotics.php content="What microcontroller would you recommend?"',$output);
$result = implode("\n",$output); $output = [];
if(strpos(strtolower($result),"pico")===false){
    echo "$R ROBOTICS EXPERT TEST FAIL3 (pi pico recommendation):$N ";
    echo "Expected to get answer containing pico but got ".$result."\n";
} else {
    echo "$G pass robotics expert test 3 (pi pico recommendation)$N\n";
}