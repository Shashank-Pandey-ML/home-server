import React from 'react';
import blog1 from '../assets/images/blog1.png';

function Blog() {
    return (
        <section className="section" id="blog">
            <div className="container text-center">
                <p className="section-subtitle">Recent Posts?</p>
                <h6 className="section-title mb-6">Blog</h6>
                <div className="blog-card">
                    <div className="blog-card-header">
                        <img src={blog1} className="blog-card-img" alt="Blog Post" />
                    </div>
                    <div className="blog-card-body">
                        <h5 className="blog-card-title">Advanced Docker Networking: Establishing Namespace Communication Across Containers</h5>
                        <p className="blog-card-caption">
                            <span>By: Shashank Pandey</span>
                        </p>
                        <p>
                            Recently, I encountered an intriguing problem where I needed to establish 
                            communication from a network namespace inside one container to another network 
                            namespace inside a different container. In this blog, I will walk you through the 
                            thought process and steps involved in solving this issue.
                        </p>
                        <p><b>Pre-requisite:</b></p>
                        <p>
                            This blog is intended for readers who are familiar with the following concepts:
                        </p>
                        <ul className="text-left">
                            <li>Docker networking</li>
                            <li>Bridge network</li>
                            <li>Virtual Ethernet (veth)</li>
                        </ul>
                        <a 
                            href="https://medium.com/@kakupandey.sp/advanced-docker-networking-establishing-namespace-communication-across-containers-56afafd15470" 
                            className="blog-card-link"
                            target="_blank"
                            rel="noopener noreferrer"
                        >
                            Read more <i className="ti-angle-double-right"></i>
                        </a>
                    </div>
                </div>
            </div>
        </section>
    );
}

export default Blog;
