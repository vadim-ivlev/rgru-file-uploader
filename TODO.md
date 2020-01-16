
Пример теста
``` javascript

function test(){
  let html =''
  if (response && !response.errors){

    let url ="http://localhost:5500/uploads"+response.data.upload_internet_file.filepath
    let color = response.data.upload_internet_file.dominant_color.hex
    html =
    `<img width=50 src="${url}">
     <div style="width: 50px; background-color:#${color}; color: silver;">...</div>` 

  } else {
    html = response.errors[0].message
  } 
  return html
}

test()

```