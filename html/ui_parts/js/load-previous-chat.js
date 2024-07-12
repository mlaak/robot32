function load_previous_chat(){
    do_load_previous_chat();
}

async function do_load_previous_chat(){
    let all_keys = Object.keys(localStorage);
    all_keys = naturalSort(all_keys);
    all_keys = all_keys.reverse();

    let displayed = 0;
    let max = 10;

    for (let i = 0; i < all_keys.length; i++) {
        if(all_keys[i].startsWith("history")){
            let r = await load_chat_part(all_keys[i]);
            if(r!=0){
                //console.log("rrrrrrrrrrrrrrrrrrrrr",r);
                displayed++;
            }
            if(displayed>=max)break; 
        }
        console.log(all_keys[i]);
    }
}

async function load_chat_part(id){
    //console.log("load chat part");
    let num = id.substring(7)*1;
    let dt = localStorage.getItem(id);

    //console.log(id,dt);
    let decrypted =await decryptText(dt,ClientSecretKey);
    //console.log("jjj",decrypted);
    let dat = JSON.parse(decrypted);

    let img_dt = localStorage.getItem("img"+num);
    //console.log("here 1");
    let img_src = "";
    if(img_dt!=null){
        img_src = await decryptText(img_dt,ClientSecretKey);
    }

    console.log("here");
    console.log(num,img_src,dat);

    return restore_interaction(num,img_src,dat);

    /*let decrypted = await window.crypto.subtle.decrypt(
        {name: "AES-GCM",iv: dt.iv},
        key,dt.encrypted
      );
      let dec = new TextDecoder();
      let txt = dec.decode(decrypted);*/
}


