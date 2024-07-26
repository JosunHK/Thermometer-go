# Thermometer-go

This is a program that tries to find the solution to the Thermometers game deterministically

for rules and testing check : https://www.puzzle-thermometers.com/

# JS script

To import custom data, open dev console in the thermometer game,
paste the code below into the console, copy and paste the result string into the `raw` variable

```javascript
let cells = [];
document.querySelectorAll(".cell-off").forEach((cell) => {
  let clsName = cell.className;
  let classes = clsName.split(" ");
  let prefix = classes[classes.length - 2] === "start" ? "s" : "~";
  cells.push(
    prefix +
      "-" +
      Math.floor(cell.style.top.split("px")[0] / 30) +
      "-" +
      Math.floor(cell.style.left.split("px")[0] / 30)
  );
});
let s1 = cells.join(",");
let vs = [];
document.querySelectorAll(".task.v").forEach((v) => {
  vs.push(v.innerText);
});
let s2 = vs.join(",");
let hs = [];
document.querySelectorAll(".task.h").forEach((h) => {
  hs.push(h.innerText);
});
let s3 = hs.join(",");
console.log(s1 + ";" + s2 + ";" + s3);
```

## NOTES

-> uncomments the `realTimePrintMap()` function to see the paths in real time
(its a significantly slower, only for debug in 4x4 is recommanded)
