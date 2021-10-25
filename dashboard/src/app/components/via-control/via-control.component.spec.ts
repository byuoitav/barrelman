import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ViaControlComponent } from './via-control.component';

describe('ViaControlComponent', () => {
  let component: ViaControlComponent;
  let fixture: ComponentFixture<ViaControlComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ViaControlComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ViaControlComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
