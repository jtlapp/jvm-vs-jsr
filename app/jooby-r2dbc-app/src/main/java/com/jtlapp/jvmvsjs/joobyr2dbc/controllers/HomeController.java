package com.jtlapp.jvmvsjs.joobyr2dbc.controllers;

import io.jooby.annotation.*;
import jakarta.inject.Singleton;

@Singleton
@Path("/")
public class HomeController {

  @GET
  public String home() {
    return "Running Jooby/MVC with R2DBC";
  }
}
