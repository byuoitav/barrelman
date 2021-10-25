import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ReachableDevicesComponent } from './reachable-devices.component';

describe('ReachableDevicesComponent', () => {
  let component: ReachableDevicesComponent;
  let fixture: ComponentFixture<ReachableDevicesComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ReachableDevicesComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ReachableDevicesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
