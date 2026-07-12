
CREATE TYPE user_role AS ENUM ('EMPLOYEE','DEPARTMENT_HEAD','ASSET_MANAGER','ADMIN');
CREATE TYPE user_status AS ENUM ('ACTIVE','INACTIVE');
CREATE TYPE department_status AS ENUM ('ACTIVE','INACTIVE');
CREATE TYPE asset_status AS ENUM ('AVAILABLE','ALLOCATED','RESERVED','UNDER_MAINTENANCE','LOST','RETIRED','DISPOSED');
CREATE TYPE maintenance_priority AS ENUM ('LOW','MEDIUM','HIGH','CRITICAL');
CREATE TYPE maintenance_status AS ENUM ('PENDING','APPROVED','REJECTED','TECHNICIAN_ASSIGNED','IN_PROGRESS','RESOLVED');
CREATE TYPE booking_status AS ENUM ('UPCOMING','ONGOING','COMPLETED','CANCELLED');
CREATE TYPE transfer_status AS ENUM ('REQUESTED','APPROVED','REJECTED');
CREATE TYPE asset_flag AS ENUM ('SHARABLE','BOOKABLE');

CREATE TABLE departments(
 id BIGSERIAL PRIMARY KEY,
 name VARCHAR(150) UNIQUE NOT NULL,
 parent_department_id BIGINT,
 head_id BIGINT,
 status department_status DEFAULT 'ACTIVE',
 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 deleted_at TIMESTAMP
);

CREATE TABLE users(
 id BIGSERIAL PRIMARY KEY,
 name VARCHAR(120) NOT NULL,
 email VARCHAR(255) UNIQUE NOT NULL,
 password TEXT NOT NULL,
 role user_role DEFAULT 'EMPLOYEE',
 gender VARCHAR(20),
 department_id BIGINT,
 status user_status DEFAULT 'ACTIVE',
 allotted_asset_id BIGINT[],
 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 deleted_at TIMESTAMP
);

CREATE TABLE assets(
 id BIGSERIAL PRIMARY KEY,
 tag VARCHAR(50) UNIQUE NOT NULL,
 name VARCHAR(150) NOT NULL,
 category VARCHAR(100) NOT NULL,
 status asset_status DEFAULT 'AVAILABLE',
 location VARCHAR(255),
 expected_location VARCHAR(255),
 condition TEXT,
 photos_docs TEXT[],
 flag asset_flag[],
 acquisition_date DATE,
 acquisition_cost NUMERIC(12,2),
 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 deleted_at TIMESTAMP
);

CREATE TABLE allocations(
 id BIGSERIAL PRIMARY KEY,
 asset_id BIGINT NOT NULL,
 from_user_id BIGINT,
 to_user_id BIGINT NOT NULL,
 allotted_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 expected_return_date TIMESTAMP,
 actual_return_date TIMESTAMP,
 reason TEXT,
 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 deleted_at TIMESTAMP
);

CREATE TABLE resource_bookings(
 id BIGSERIAL PRIMARY KEY,
 asset_id BIGINT NOT NULL,
 booked_by BIGINT NOT NULL,
 booking_status booking_status DEFAULT 'UPCOMING',
 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 deleted_at TIMESTAMP
);

CREATE TABLE booking_slots(
 id BIGSERIAL PRIMARY KEY,
 booking_id BIGINT NOT NULL,
 start_time TIMESTAMP NOT NULL,
 end_time TIMESTAMP NOT NULL,
 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 deleted_at TIMESTAMP
);

CREATE TABLE maintenance(
 id BIGSERIAL PRIMARY KEY,
 asset_id BIGINT NOT NULL,
 priority maintenance_priority,
 description TEXT,
 images TEXT[],
 status maintenance_status DEFAULT 'PENDING',
 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 deleted_at TIMESTAMP
);

CREATE TABLE maintenance_workflows(
 id BIGSERIAL PRIMARY KEY,
 maintenance_id BIGINT NOT NULL,
 status maintenance_status NOT NULL,
 description TEXT,
 updated_by BIGINT,
 workflow_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 deleted_at TIMESTAMP
);

CREATE TABLE audits(
 id BIGSERIAL PRIMARY KEY,
 auditors BIGINT[],
 scope TEXT,
 from_date DATE,
 to_date DATE,
 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 deleted_at TIMESTAMP
);

CREATE TABLE audit_reports(
 id BIGSERIAL PRIMARY KEY,
 audit_id BIGINT NOT NULL,
 asset_id BIGINT NOT NULL,
 verified BOOLEAN,
 status VARCHAR(50),
 remarks TEXT,
 verified_by BIGINT,
 verified_at TIMESTAMP,
 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 deleted_at TIMESTAMP
);

CREATE TABLE transfer_requests(
 id BIGSERIAL PRIMARY KEY,
 asset_id BIGINT NOT NULL,
 from_user_id BIGINT NOT NULL,
 to_user_id BIGINT NOT NULL,
 requested_by BIGINT NOT NULL,
 approved_by BIGINT,
 status transfer_status DEFAULT 'REQUESTED',
 reason TEXT,
 requested_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 approved_at TIMESTAMP,
 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
 deleted_at TIMESTAMP
);

ALTER TABLE departments ADD FOREIGN KEY(parent_department_id) REFERENCES departments(id);
ALTER TABLE departments ADD FOREIGN KEY(head_id) REFERENCES users(id);
ALTER TABLE users ADD FOREIGN KEY(department_id) REFERENCES departments(id);
ALTER TABLE allocations ADD FOREIGN KEY(asset_id) REFERENCES assets(id);
ALTER TABLE allocations ADD FOREIGN KEY(from_user_id) REFERENCES users(id);
ALTER TABLE allocations ADD FOREIGN KEY(to_user_id) REFERENCES users(id);
ALTER TABLE resource_bookings ADD FOREIGN KEY(asset_id) REFERENCES assets(id);
ALTER TABLE resource_bookings ADD FOREIGN KEY(booked_by) REFERENCES users(id);
ALTER TABLE booking_slots ADD FOREIGN KEY(booking_id) REFERENCES resource_bookings(id) ON DELETE CASCADE;
ALTER TABLE maintenance ADD FOREIGN KEY(asset_id) REFERENCES assets(id);
ALTER TABLE maintenance_workflows ADD FOREIGN KEY(maintenance_id) REFERENCES maintenance(id) ON DELETE CASCADE;
ALTER TABLE maintenance_workflows ADD FOREIGN KEY(updated_by) REFERENCES users(id);
ALTER TABLE audit_reports ADD FOREIGN KEY(audit_id) REFERENCES audits(id) ON DELETE CASCADE;
ALTER TABLE audit_reports ADD FOREIGN KEY(asset_id) REFERENCES assets(id);
ALTER TABLE audit_reports ADD FOREIGN KEY(verified_by) REFERENCES users(id);
ALTER TABLE transfer_requests ADD FOREIGN KEY(asset_id) REFERENCES assets(id);
ALTER TABLE transfer_requests ADD FOREIGN KEY(from_user_id) REFERENCES users(id);
ALTER TABLE transfer_requests ADD FOREIGN KEY(to_user_id) REFERENCES users(id);
ALTER TABLE transfer_requests ADD FOREIGN KEY(requested_by) REFERENCES users(id);
ALTER TABLE transfer_requests ADD FOREIGN KEY(approved_by) REFERENCES users(id);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_assets_tag ON assets(tag);
CREATE INDEX idx_assets_status ON assets(status);
CREATE INDEX idx_alloc_asset ON allocations(asset_id);
CREATE INDEX idx_booking_asset ON resource_bookings(asset_id);
CREATE INDEX idx_booking_slot_time ON booking_slots(start_time,end_time);
CREATE INDEX idx_maint_asset ON maintenance(asset_id);
CREATE INDEX idx_transfer_asset ON transfer_requests(asset_id);
CREATE INDEX idx_audit_asset ON audit_reports(asset_id);
