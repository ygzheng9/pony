const Glob = require("glob");

function getSources() {
    var entries = {
      application: [
        './node_modules/jquery-ujs/src/rails.js',
        './assets/css/application.scss',
      ],
    }


    // all files in ./assets folder  
    // Glob.sync("./assets/*/*.*").forEach((entry) => {
    Glob.sync("./assets/{js,css}/**/*.*").forEach((entry) => {
      if (entry === './assets/css/application.scss') {
        return
      }

      // ignore those folders  
      // if ((/(\/vendors\/|\/src\/|\/data\/)/i).test(entry)) {
      //   // console.log("got it: ", entry)
      //   return 
      // }

      // only those file types 
      if ( (/(ts|js|s[ac]ss|go)$/i).test(entry) == false ) {
        return 
      }

      // ignore files name start with underscore
      let key = entry.replace(/(\.\/assets\/(src|js|css|go)\/)|\.(ts|js|s[ac]ss|go)/g, '')
      if(key.startsWith("_")) {
        return
      }

      if( entries[key] == null) {
        entries[key] = [entry]
        return
      }

      entries[key].push(entry)
    })
    return entries
}

var a = getSources(); 

console.log("\n\nfinal result."); 
console.log(a);

// for (let v of a) {
//   console.log(v[0], "-", v[1]); 
// }

// a.forEach((k, v) => console.log(k, v));
