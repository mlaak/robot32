<?php
require __DIR__."/vendor/Robot32lib/Middleware/Middleware.php";

$FAL_KEY = trim(file_get_contents(__DIR__."/../../../keys/falai.txt"));
$URL = "https://fal.run/fal-ai/fast-turbo-diffusion";
$IMAGES_DIR = "../../recieved_images";


$content = $_REQUEST["content"];
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

$result = curl_exec($ch);
curl_close($ch);

$json = json_decode($result ,true);
$dat = base64_decode(explode(',',$json['images'][0]['url'],2)[1]);
$md5 = md5($dat);
$time = microtime(true);
$fname = "$IMAGES_DIR/$md5.$time.jpg";

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









