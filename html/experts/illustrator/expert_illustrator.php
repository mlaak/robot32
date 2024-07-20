<?php
require __DIR__."/settings.php";
require __DIR__."/vendor/Robot32lib/Middleware/Middleware.php";
include __DIR__."/vendor/Robot32lib/ImageSource/ImageSource.php";
use Robot32lib\ImageSource\ImageSource; 



$use_falai = true;
$falai_error = false;

$content = $_REQUEST["content"];

if(isset($_REQUEST["prerendered"]))$use_falai = false;
if(!$FAL_KEY && trim($FAL_KEY)==="")$use_falai = false; 

if($use_falai){
    $data = ["prompt" => $content];    
    $ch = curl_init();
    curl_setopt($ch, CURLOPT_URL, $URL);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($ch, CURLOPT_POST, true);
    curl_setopt($ch, CURLOPT_HTTPHEADER, array(
        "Authorization: Key $FAL_KEY",
        "Content-type: application/json"
    ));
    curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($data));

    curl_setopt($ch, CURLOPT_CONNECTTIMEOUT, 15); 
    curl_setopt($ch, CURLOPT_TIMEOUT, 15);

    $result = curl_exec($ch);
    curl_close($ch);

    $json = json_decode($result ,true);
    //print_r($json);
    if($json===null || !isset($json['images'][0]['url'])){
        $falai_error = true;
        error_log("ERROR GETTING IMAGE FROM FALAI!");
    }
    else $dat = base64_decode(explode(',',$json['images'][0]['url'],2)[1]);
}
if(!$use_falai || $falai_error){
    if(class_exists("Robot32lib\\ImageSource\\ImageSource")){
        $is = new ImageSource(); 
        if(strpos($content,"robot")!==false)
            $dat=$is->getRandomPicture(["robot"]);
        else if(strpos($content,"PHP")!==false)
            $dat=$is->getRandomPicture(["PHP"]);
        else if(strpos($content,"golang")!==false)
            $dat=$is->getRandomPicture(["golang"]);       
        else if(strpos($content,"javascript")!==false)
            $dat=$is->getRandomPicture(["javascript"]);
        else if(strpos($content,"electronics")!==false)
            $dat=$is->getRandomPicture(["electronics"]);
        else if(strpos($content,"joke")!==false)
            $dat=$is->getRandomPicture(["joke"]);
        else if(strpos($content,"microcontroller")!==false)
            $dat=$is->getRandomPicture(["microcontroller"]);
        else $dat=$is->getRandomPicture(["general"]);

        header("Prerendered-Image:1");
    }   
}



$md5 = md5($dat);
$time = str_replace(".","_",microtime(true));

//so in the future the reverse proxy could save the image instead
header("Return-Image:".base64_encode($dat));
header("Return-Image-Name:".base64_encode($md5."_".$time));


$fname = "$IMAGES_DIR/$md5"."_"."$time.jpg";

if(!isset($_REQUEST['raw'])){
    file_put_contents(__DIR__."/".$fname, $dat);
    echo "{\"image\":\"$fname\"}";
}
else {
    if($_REQUEST['raw']=="base64"){
        echo base64_encode($dat);
    }
    else echo $dat;
}









