console.log("Hello code!");

let resizeExploreElm = document.getElementById("resizeExplorer");

// The current position of mouse
let x = 0;
let y = 0;

// The dimension of the element
let w = 0;
let h = 0;

// Handle the mousedown event
// that's triggered when user drags the resizer
const mouseDownHandler = function (e) {

    //let ele = e.target;
    //console.log(ele);

    // Get the current mouse position
    x = e.clientX;
    y = e.clientY;

    // Calculate the dimension of element
    const styles = window.getComputedStyle(resizeExploreElm);
    w = parseInt(styles.width, 10);
    h = parseInt(styles.height, 10);

    // Attach the listeners to `document`
    document.addEventListener('mousemove', mouseMoveHandler);
    document.addEventListener('mouseup', mouseUpHandler);
};


const mouseMoveHandler = function (e) {

    // How far the mouse has been moved
    const dx = e.clientX - x;
    const dy = e.clientY - y;

    //console.log("resizeExploreElm:",resizeExploreElm);
    //console.log("target:",e.target);

    // Adjust the dimension of element
    resizeExploreElm.style.width = `${w + dx}px`;
    //resizeExploreElm.style.height = `${h + dy}px`;
    //console.log(`${w + dx}px`);
};

const mouseUpHandler = function () {
    // Remove the handlers of `mousemove` and `mouseup`
    document.removeEventListener('mousemove', mouseMoveHandler);
    document.removeEventListener('mouseup', mouseUpHandler);
};


// Query all resizers
const resizers = resizeExploreElm.querySelectorAll('.resizer');

// Loop over them
[].forEach.call(resizers, function (resizer) {
    resizer.addEventListener('mousedown', mouseDownHandler);
});
