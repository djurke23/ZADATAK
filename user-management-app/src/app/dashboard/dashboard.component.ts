import { Component, OnInit } from '@angular/core';
import { AuthService } from '../services/auth.service';
import { CommonModule } from '@angular/common'; // Uvoz CommonModule

@Component({
  standalone: true,
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css'],
  imports: [CommonModule] // Dodaj CommonModule ovde
})
export class DashboardComponent implements OnInit {
  users: any[] = [];

  constructor(private authService: AuthService) {}

  ngOnInit() {
    this.loadUsers();
  }

  loadUsers() {
    this.authService.getUsers().subscribe({
      next: (response: any[]) => {
        this.users = response;
      },
      error: (err: any) => {
        alert('Failed to load users: ' + (err.error?.message || 'Unknown error'));
      }
    });
  }
}
