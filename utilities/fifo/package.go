// Package fifo provides FIFO (First In, First Out) fixed-size storage structures. This is useful
// when determining whether an item has already been seen within a certain period.
//
// # Contains
//
//   - Store - FIFO fixed-size storage with fast-lookup.
//   - List - Fixed-size doubly-linked list. (Used by Store for tracking order).
//   - Node - Item type by both Store and List.
//   - Simple - An extremely simplistic FIFO store implementation.
//
// # Scenario
//
// Live telemetry files contain 3 data buffers from which telemetry ticks are parsed. Each buffer will
// specify a tick number (this can be though of as an ID) to avoid any duplication. However, we may incur the
// the same tick number multiple times when trying to read these buffers. Thus, the objectives of the FIFO package is
// to provide structures that can tell us whether a tick was seen recently and remove ticks that have not been seen
// recently.
package fifo
