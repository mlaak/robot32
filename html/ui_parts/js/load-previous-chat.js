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
                displayed++;
            }
            if(displayed>=max)break; 
        }
        console.log(all_keys[i]);
    }
}

async function load_chat_part(id){
    let num = id.substring(7)*1;
    let dt = localStorage.getItem(id);

    let decrypted =await decryptText(dt,ClientSecretKey);
    let dat = JSON.parse(decrypted);

    let img_dt = localStorage.getItem("img"+num);
    let img_src = "";
    if(img_dt!=null){
        img_src = await decryptText(img_dt,ClientSecretKey);
    }

    return restore_interaction(num,img_src,dat);
}


