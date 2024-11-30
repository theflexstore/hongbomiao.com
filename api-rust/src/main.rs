use axum::{routing::get, serve, Router};
use http::Method;
use tokio::net::TcpListener;
use tower_http::cors::{Any, CorsLayer};
use tracing::info;
use std::env;

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt()
        .with_max_level(tracing::Level::TRACE)
        .init();

    #[cfg(debug_assertions)]
    dotenvy::from_filename(".env.development").ok();
    #[cfg(not(debug_assertions))]
    dotenvy::from_filename(".env.production").ok();

    // Read port from environment variable, with fallback to 3000
    let port = env::var("PORT")
        .expect("PORT must be set in environment")
        .parse::<u16>()
        .expect("PORT must be a valid number");

    // Create CORS middleware
    let cors = CorsLayer::new()
        .allow_methods([Method::GET, Method::POST])
        .allow_origin(Any)
        .allow_headers(Any);

    let app = Router::new()
        .route("/", get(root))
        .layer(cors);

    let listener = TcpListener::bind(format!("0.0.0.0:{}", port)).await.unwrap();
    info!("Server listening on port {}", port);

    serve(listener, app).await.unwrap();
}

async fn root() -> &'static str {
    "Hello, World!"
}
