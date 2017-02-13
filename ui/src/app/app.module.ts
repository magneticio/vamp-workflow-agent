import {NgModule} from '@angular/core';
import {MaterialModule} from '@angular/material';
import {BrowserModule} from '@angular/platform-browser';

import 'hammerjs';
import {AppComponent} from './app.component';

@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    BrowserModule,
    MaterialModule.forRoot(),
  ],
  providers: [],
  entryComponents: [],
  bootstrap: [AppComponent]
})
export class AppModule {
}
