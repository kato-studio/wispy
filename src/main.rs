use actix_web::{get, post, web, App, HttpResponse, HttpServer, Responder};

#[get("/")]
async fn index() -> impl Responder {
    "Hello, World!"
}

#[get("/html")]
async fn route_html(name: web::Path<String>) -> impl Responder {
    // return html file at relative path
    HttpResponse::Ok().body(include_str!("../view/test.kato"))
}

#[get("/_n/{name}")]
async fn hello(name: web::Path<String>) -> impl Responder {
    let messages = vec![String::from("Message 1"), String::from("<Message 2>")];

    // print current directly folders
    println!("{:?}", std::env::current_dir().unwrap());

    HttpResponse::Ok().body(format!("Hello, {}!", name))
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    HttpServer::new(|| App::new().service(index).service(hello))
        .bind(("127.0.0.1", 8080))?
        .run()
        .await
}
