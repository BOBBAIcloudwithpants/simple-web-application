// $("btn").click(function(e) {
//     console.log("aaaa")
//     e.preventDefault();
//     $.ajax({
//         type: "POST",
//         url: "/login",
//         data: {
//             username: $("username").val(),
//             access_token: $("password").val()
//         },
//         success: function(result) {
//             alert('ok');
//         },
//         error: function(result) {
//             alert('error');
//         }
//     });
// });
function getRandomInt(max) {
    return Math.floor(Math.random() * Math.floor(max));
}
function login(username, password) {
    console.log(username)
    console.log(password)
    let randomid = getRandomInt(1000000)
    $.ajax({
        type: "POST",
        url: "/",
        data: {
            username: username,
            password: password,
            user_id: randomid
        },
        success: function(result) {
            alert('登录成功');
            window.location.href = `http://localhost:3000/users/${randomid}`;
        },
        error: function(result) {
            alert('error');
        }
    });
}